package ru.kolobanov.pc.club.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.http.*;
import ru.kolobanov.pc.club.controllers.dto.*;
import ru.kolobanov.pc.club.entity.User;
import ru.kolobanov.pc.club.exeptions.*;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.repository.UserRepo;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;


@Service
@AllArgsConstructor
public class UserService {

    private UserRepo userRepo;
    private PaymentService paymentService;

    public Page<ResponseUser> getUsers(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAll(pageable).map(DtoMapper::userToResponseUser);
    }

    public Page<ResponseUserTopByBalance> getUsersOrderByBalance(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByBalance(pageable)
                .map(DtoMapper::userToResponseTopByBalance);
    }

    public Page<ResponseUserTopBySessions> getUsersOrderByCountSessions(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderBySessionsCount(pageable)
                .map(DtoMapper::userToResponseTopBySessions);
    }

    public Page<ResponseUserTopByHours> getUsersOrderByHours(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByTotalHours(pageable)
                .map(DtoMapper::userToResponseTopByHours);
    }

    public ResponseUser getUserById(Long id){
        User user =  checkUserId(id);
        return DtoMapper.userToResponseUser(user);
    }

    @Transactional
    public ResponseUser saveRegistrationUser(RequestRegistrationUser requestRegistrationUser){

        checkEmail(requestRegistrationUser.getEmail(), 0L);
        checkPhone(requestRegistrationUser.getPhone(),0L);

        User user = DtoMapper.requestRegistrationUserToUser(requestRegistrationUser);

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);


    }


    @Transactional
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
    public ResponseUser update(RequestUpdateUser requestUpdateUser){
        User user = checkUserId(requestUpdateUser.getId());

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
    public ResponseUser addBalance(RequestTopUpUserBalance requestTopUpUserBalance){
        User user = checkUserId(requestTopUpUserBalance.getUser_id());

        paymentService.receivingPayment(requestTopUpUserBalance);

        user.setBalance(user.getBalance() + requestTopUpUserBalance.getAmount());

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);

    }

    @Transactional
    public void updateBalance(Long id,Double balance){
        User user = checkUserId(id);

        user.setBalance(balance);

        userRepo.save(user);

    }

    @Transactional
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



}
