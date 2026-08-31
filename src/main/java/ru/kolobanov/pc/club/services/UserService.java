package ru.kolobanov.pc.club.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.cache.annotation.Caching;
import org.springframework.http.*;
import ru.kolobanov.pc.club.controllers.dto.*;
import ru.kolobanov.pc.club.entity.Referrals;
import ru.kolobanov.pc.club.entity.User;
import ru.kolobanov.pc.club.exeptions.*;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.repository.ReferralsRepo;
import ru.kolobanov.pc.club.repository.UserRepo;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Random;


@Service
@AllArgsConstructor
public class UserService {

    private UserRepo userRepo;
    private ReferralsRepo referralsRepo;

    @Cacheable(value = "users", key = "#page + '-' + #size")
    public List<ResponseUser> getUsers(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAll(pageable).map(DtoMapper::userToResponseUser).getContent();
    }

    @Cacheable(value = "topByBalance", key = "#page + '-' + #size")
    public List<ResponseUserTopByBalance> getUsersOrderByBalance(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByBalance(pageable)
                .map(DtoMapper::userToResponseTopByBalance).getContent();
    }

    @Cacheable(value = "topByCountSessions", key = "#page + '-' + #size")
    public List<ResponseUserTopBySessions> getUsersOrderByCountSessions(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderBySessionsCount(pageable)
                .map(DtoMapper::userToResponseTopBySessions).getContent();
    }

    @Cacheable(value = "topByHours", key = "#page + '-' + #size")
    public List<ResponseUserTopByHours> getUsersOrderByHours(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByTotalHours(pageable)
                .map(DtoMapper::userToResponseTopByHours).getContent();
    }

    @Cacheable(value = "user", key = "#id")
    public ResponseUser getUserById(Long id){
        User user =  checkUserId(id);
        return DtoMapper.userToResponseUser(user);
    }

    @Transactional
    @CacheEvict(value = {"users", "topByHours","topByCountSessions","topByBalance"}, allEntries = true)
    public ResponseUser saveRegistrationUser(RequestRegistrationUser requestRegistrationUser, Long token){

        checkEmail(requestRegistrationUser.getEmail(), 0L);
        checkPhone(requestRegistrationUser.getPhone(),0L);

        User user = DtoMapper.requestRegistrationUserToUser(requestRegistrationUser);

        user.setToken(new Random().nextLong(-1000,1001));

        userRepo.save(user);

        Long idSender;

        if(token != null){

            idSender = getUserIdByToken(token);
            Referrals referrals = DtoMapper.referralsRequestToReferral(idSender,user.getId());
            referralsRepo.save(referrals);

        }


        return DtoMapper.userToResponseUser(user);


    }


    @Transactional
    @CacheEvict(value = {"users", "topByHours","topByCountSessions","topByBalance","user"}, allEntries = true)
    public boolean deleteUserById(Long id){
        checkUserId(id);
        userRepo.deleteById(id);
        return true;
    }


    public ResponseUser loginUser(RequestLoginUser requestLoginUser){

        String email = requestLoginUser.getEmail();
        String password = requestLoginUser.getPassword();

        User user = userRepo.findUserByEmail(email);

        if(user == null){
            throw new InvalidCredentialsException("Пользователя с такой почтой не существует");
        }

        if(!password.equals(user.getPassword())){
            throw new InvalidCredentialsException("Неверный пароль");
        }

        return DtoMapper.userToResponseUser(user);
    }


    @Transactional
    @Caching(evict = {
            @CacheEvict( value = {"users","topByHours","topByCountSessions","topByBalance"}, allEntries = true),
            @CacheEvict(value = "user", key = "#id")
    })
    public ResponseUser update(Long id,RequestUpdateUser requestUpdateUser){
        User user = checkUserId(id);

        String email = requestUpdateUser.getEmail();
        String name = requestUpdateUser.getName();
        String phone = requestUpdateUser.getPhone();
        String password = requestUpdateUser.getPassword();

        if(name != null && (!name.isBlank())){
            user.setName(name);
        }

        if(email != null && (!email.isBlank())){
            checkEmail(email,user.getId());
            user.setEmail(email);
        }

        if(phone != null && (!phone.isBlank())){
            checkPhone(phone,user.getId());
            user.setNumber(phone);
        }

        if(password != null && (!password.isBlank())){
            user.setPassword(password);
        }

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);
    }

    @Transactional
    @Caching(evict = {
            @CacheEvict( value = {"users","topByBalance"}, allEntries = true),
            @CacheEvict(value = "user", key = "#id")
    })
    public ResponseUser addBalance(Long id,Double amount){
        User user = checkUserId(id);

        user.setBalance(user.getBalance() + amount);

        List<Referrals> referrals = referralsRepo.findAllByIdRef(user.getId());

        if(referrals != null){
            for(Referrals r: referrals){
                User senderRef = checkUserId(r.getIdSender());
                senderRef.setBalance(senderRef.getBalance() + amount*0.05);

                userRepo.save(senderRef);

            }
        }

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);

    }

    @Transactional
    @Caching(evict = {
            @CacheEvict( value = {"users","topByBalance"}, allEntries = true),
            @CacheEvict(value = "user", key = "#id")
    })
    public void updateBalance(Long id,Double balance){
        User user = checkUserId(id);

        user.setBalance(balance);

        userRepo.save(user);

    }

    @Transactional
    @Caching(evict = {
            @CacheEvict( value = {"users","topByHours"}, allEntries = true),
            @CacheEvict(value = "user", key = "#id")
    })
    public void updateHours(Long id, Double hours){
        User user = checkUserId(id);

        user.setHours(hours);

        userRepo.save(user);
    }

    public User checkUserId(Long id){
        return userRepo.findById(id)
                .orElseThrow(() -> new UserNotFoundException("Пользователя с таким Id нету"));
    }

    private void checkEmail(String email,Long id){
        User user = userRepo.findUserByEmail(email);

        if(user != null && !(user.getId().equals(id))){
            throw new EmailAlreadyExistException("Пользователь с данной почтой уже существует");
        }
    }

    private void checkPhone(String phone,Long id){

        User user = userRepo.findUserByNumber(phone);

        if(user != null && !(user.getId().equals(id))){
            throw new PhoneAlreadyExistException("Пользователь с таким номером уже существует");

        }
    }

    public boolean chekUserBalance(User user,Double total){
        if(user.getBalance() < total){
            throw new InsufficientBalanceException("Пополните баланс");
        }

        return true;
    }

    public Long getUserIdByToken(Long token){
        User user = userRepo.findUserByToken(token);
        if(user != null){
            return user.getId();
        }
        return null;
    }



}
