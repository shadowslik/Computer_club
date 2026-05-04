package org.example.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.example.controllers.dto.RequestUser;
import org.example.controllers.dto.ResponseUser;
import org.example.entity.User;
import org.example.mapper.DtoMapper;
import org.example.repository.UserRepo;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
@AllArgsConstructor
public class UserService {

    private UserRepo userRepo;

    public Page<ResponseUser> getUsers(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAll(pageable).map(DtoMapper::userToResponseUser);
    }

    public Page<ResponseUser> getUsersOrderByBalance(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByBalance(pageable)
                .map(DtoMapper::userToResponseUser);
    }

    public Page<ResponseUser> getUsersOrderByCountSessions(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByComputerSessionsCount(pageable)
                .map(DtoMapper::userToResponseUser);
    }

    public Page<ResponseUser> getUsersOrderByHours(int page, int size){
        Pageable pageable = PageRequest.of(page,size);
        return userRepo.findAllOrderByTotalHours(pageable)
                .map(DtoMapper::userToResponseUser);
    }

    public ResponseUser getUserById(Long id){
        User user =  checkUserId(id);
        return DtoMapper.userToResponseUser(user);
    }

    @Transactional
    public ResponseUser saveRegistrationUser(RequestUser requestUser){

        checkPostParams(requestUser);

        User user = DtoMapper.requestUserToUser(requestUser);

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);


    }


    @Transactional
    public void deleteUserById(Long id){

        userRepo.findById(id).orElseThrow(() -> new IllegalArgumentException("Такого id нету в базе данных"));
        userRepo.deleteById(id);
    }


    public ResponseUser loginUser(String email,String password){
        User user = userRepo.findUserByEmail(email);
        if(!checkEmail(email)){
            throw new IllegalArgumentException("Email не корректный");
        }
        if(user == null){
            throw new IllegalArgumentException("Пользователя с такой почтой не существует");
        }

        if(!user.getPassword().equals(password)){
            throw new IllegalArgumentException("Неверный пароль");
        }

        return DtoMapper.userToResponseUser(user);
    }


    public ResponseUser update(Long id,RequestUser requestUser){
        User user = userRepo.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Пользовател с таким id нет"));



        System.out.println(requestUser.toString());

        checkUpdateParams(requestUser);

        if(!requestUser.getName().equals("null")){
            user.setName(requestUser.getName());
        }

        if(!requestUser.getEmail().equals("null")){
            user.setEmail(requestUser.getEmail());
        }

        if(!(requestUser.getPassword() == null)){
            user.setPassword(requestUser.getPassword());
        }

        if(!requestUser.getPhone().equals("null")){
            user.setNumber(requestUser.getPhone());
        }

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);
    }

    public ResponseUser updateBalance(Long id,Double balance){
        User user = userRepo.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Пользовател с таким id нет"));

        if(balance < 0){
            throw new IllegalArgumentException("Баланс должен быть положительным");
        }
        user.setBalance(balance);

        userRepo.save(user);

        return DtoMapper.userToResponseUser(user);

    }


    private boolean checkEmail(String email){
        String regex = "^[\\w.-]+@[\\w.-]+\\.[a-zA-Z]{2,}$";
        return email.matches(regex);
    }

    private boolean checkPhone(String phone){
        String regex = "^[+\\d][\\d\\s\\-()]{7,20}$";
        return phone.matches(regex);
    }

    public User checkUserId(Long id){
        return userRepo.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Пользователя с таким Id нету"));
    }

    private void checkPostParams(RequestUser requestUser){

        checkUpdateParams(requestUser);

        if(!requestUser.getEmail().equals("null") && userRepo.findUserByEmail(requestUser.getEmail()) != null){
            throw new IllegalArgumentException("Пользователь с данной почтой уже существует");
        }

        if(!requestUser.getPhone().equals("null") && userRepo.findUserByNumber(requestUser.getPhone()) != null){
            throw new IllegalArgumentException("Пользователь с таким номером уже существует");

        }
    }

    private void checkUpdateParams(RequestUser requestUser){
        if(!requestUser.getName().equals("null") && requestUser.getName().length() < 2){
            throw new IllegalArgumentException("Имя слишком короткое");
        }
        if(!requestUser.getEmail().equals("null") && !checkEmail(requestUser.getEmail())){
            throw new IllegalArgumentException("Email не корректный");
        }

        if(!requestUser.getPhone().equals("null") && !checkPhone(requestUser.getPhone())){
            throw new IllegalArgumentException("Телефон не корректный");
        }

        if(!(requestUser.getPassword() == null) && requestUser.getPassword().length() < 6){
            throw new IllegalArgumentException("Пароль слишком короткий");
        }
    }


}
