package org.example.controllers;

import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.AllArgsConstructor;
import org.example.controllers.dto.RequestUser;
import org.example.controllers.dto.ResponseUser;
import org.example.mapper.DtoMapper;
import org.example.services.UserService;
import org.springframework.data.domain.Page;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;


@RestController
@AllArgsConstructor
@RequestMapping("api/clients")
@CrossOrigin(origins = "http://localhost:8081")
@Tag(name = "Управление клиентами")
public class UserController {

    private UserService userService;


    @GetMapping
    public ResponseEntity<List<ResponseUser>> getAllUsers(@RequestParam int page,
                                                          @RequestParam int size){

        return ResponseEntity.ok(userService.getUsers(page,size)
                .getContent());
    }

    @GetMapping("/{id}")
    public ResponseEntity<ResponseUser> getUserById(@PathVariable Long id){
        return ResponseEntity.ok(userService.getUserById(id));
    }

    @GetMapping("/sort_by_balance")
    public ResponseEntity<List<ResponseUser>> getAllUsersOrderByBalance(@RequestParam int page,
                                                                        @RequestParam int size){
        return ResponseEntity.ok(userService.getUsersOrderByBalance(page,size)
                .getContent());
    }

    @GetMapping("/sort_by_count_sessions")
    public ResponseEntity<List<ResponseUser>> getAllUsersOrderByCountSessions(@RequestParam int page,
                                                                        @RequestParam int size){
        return ResponseEntity.ok(userService.getUsersOrderByCountSessions(page,size)
                .getContent());
    }

    @GetMapping("/sort_by_count_hours")
    public ResponseEntity<List<ResponseUser>> getAllUsersOrderByHours(@RequestParam int page,
                                                                              @RequestParam int size){
        return ResponseEntity.ok(userService.getUsersOrderByHours(page,size)
                .getContent());
    }

    @PostMapping
    public ResponseEntity<ResponseUser> registrationUser(@RequestBody RequestUser requestUser){
        return ResponseEntity.status(HttpStatus.CREATED).body(userService.saveRegistrationUser(requestUser));
    }

    @PostMapping("/login")
    public ResponseEntity<ResponseUser> loginUser(
            @RequestParam String email,
            @RequestParam String password
    ){
        return ResponseEntity.status(HttpStatus.CREATED).body(userService.loginUser(email, password));
    }

    @PutMapping
    public ResponseEntity<ResponseUser> updateUser(
            @RequestParam Long id,
            @RequestBody RequestUser requestUser){
        return ResponseEntity.status(HttpStatus.OK).body(userService.update(id,requestUser));
    }

    @PutMapping("/{balance}")
    public ResponseEntity<ResponseUser> updateBalance(
            @RequestParam Long id,
            @RequestParam Double balance
    ){
        return ResponseEntity.status(HttpStatus.OK).body(userService.updateBalance(id,balance));

    }


    @DeleteMapping
    public void deleteUser(@RequestParam Long id){
        userService.deleteUserById(id);
    }
}
