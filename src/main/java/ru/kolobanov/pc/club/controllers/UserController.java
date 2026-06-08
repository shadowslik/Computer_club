package ru.kolobanov.pc.club.controllers;

import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import org.springframework.validation.annotation.Validated;
import ru.kolobanov.pc.club.controllers.dto.*;
import ru.kolobanov.pc.club.services.UserService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;


@RestController
@AllArgsConstructor
@RequestMapping("api/users")
@Tag(name = "Управление клиентами")
@Validated
public class UserController {

    private UserService userService;


    @GetMapping
    public ResponseEntity<List<ResponseUser>> getUsers(@RequestParam int page,
                                                          @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(userService.getUsers(page,size).getContent());
    }

    @GetMapping("/id")
    public ResponseEntity<ResponseUser> getUserById(
            @RequestParam
            @Min(value = 1,message = "id пользователя должен быть больше 0")
            Long id){
        return ResponseEntity.status(HttpStatus.OK).body(userService.getUserById(id));
    }

    @GetMapping("/sort_by_balance")
    public ResponseEntity<List<ResponseUserTopByBalance>> getAllUsersOrderByBalance(@RequestParam int page,
                                                                                    @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(userService.getUsersOrderByBalance(page,size)
                .getContent());
    }

    @GetMapping("/sort_by_count_sessions")
    public ResponseEntity<List<ResponseUserTopBySessions>> getAllUsersOrderByCountSessions(@RequestParam int page,
                                                                                           @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(userService.getUsersOrderByCountSessions(page,size)
                .getContent());
    }

    @GetMapping("/sort_by_count_hours")
    public ResponseEntity<List<ResponseUserTopByHours>> getAllUsersOrderByHours(@RequestParam int page,
                                                                                @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(userService.getUsersOrderByHours(page,size)
                .getContent());
    }

    @PostMapping
    public ResponseEntity<ResponseUser> registrationUser(@Valid @RequestBody RequestRegistrationUser requestRegistrationUser){
        return ResponseEntity.status(HttpStatus.CREATED).body(userService.saveRegistrationUser(requestRegistrationUser));
    }

    @PostMapping("/login")
    public ResponseEntity<ResponseUser> loginUser(@Valid @RequestBody RequestLoginUser requestLoginUser){
        return ResponseEntity.status(HttpStatus.OK).body(userService.loginUser(requestLoginUser));
    }

    @PatchMapping
    public ResponseEntity<ResponseUser> updateUser(@Valid @RequestBody RequestUpdateUser requestUpdateUser){
        return ResponseEntity.status(HttpStatus.OK).body(userService.update(requestUpdateUser));
    }

    @PatchMapping("/balance")
    public ResponseEntity<ResponseUser> updateBalance(@Valid @RequestBody RequestTopUpUserBalance requestTopUpUserBalance){
        return ResponseEntity.status(HttpStatus.OK).body(userService.addBalance(requestTopUpUserBalance));

    }

    @DeleteMapping
    public ResponseEntity<Boolean> deleteUser(
            @RequestParam
            @Min(value = 1,message = "id пользователя должен быть больше 0")
            Long id){
        return ResponseEntity.status(HttpStatus.OK).body(userService.deleteUserById(id));
    }
}
