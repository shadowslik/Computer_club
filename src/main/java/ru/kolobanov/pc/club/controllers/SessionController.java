package ru.kolobanov.pc.club.controllers;


import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Min;
import lombok.AllArgsConstructor;
import org.springframework.validation.annotation.Validated;
import ru.kolobanov.pc.club.controllers.dto.*;
import ru.kolobanov.pc.club.services.SessionService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@AllArgsConstructor
@RequestMapping("api/sessions")
@Tag(name = "Управление сессиями")
@Validated
public class SessionController {

    private SessionService sessionService;

    @GetMapping
    public ResponseEntity<List<ResponseSessionAdmin>> findAll(@RequestParam int page,
                                                         @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.getSessions(page,size).getContent());
    }

    @PostMapping
    public ResponseEntity<ResponseSession> post(@Valid @RequestBody RequestCreateSession requestCreateSession){
        return ResponseEntity.status(HttpStatus.CREATED).body(sessionService.setComputerSession(requestCreateSession));
    }

    @GetMapping("/id")
    public ResponseEntity<ResponseSessionAdmin> getById(@RequestParam @Min(value = 1,message = "id сессии должен быть больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.getSession(id));
    }

    @GetMapping("/lastUserSessions")
    public ResponseEntity<List<ResponseLastUserSessions>> getByUserId(@RequestParam @Min(value = 1,message = "id пользователя должен быть больше 0") Long user_id){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.getLastUserSessions(user_id));
    }

    @GetMapping("/bookedUserSessions")
    public ResponseEntity<List<ResponseSession>> getBookedByUserId(@RequestParam @Min(value = 1,message = "id пользователя должен быть больше 0") Long user_id,
                                                                   @RequestParam int page,
                                                                   @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.getSessionsByUserId(user_id,page,size)
                .getContent());
    }

    @GetMapping("/computerSessions")
    public ResponseEntity<List<ResponseSessionAdmin>> getByComputerId(@RequestParam @Min(value = 1,message = "id компьютера должен быть больше 0") Long computer_id,
                                                                 @RequestParam int page,
                                                                 @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService
                .getSessionsByComputerId(computer_id,page,size)
                .getContent());
    }

    @DeleteMapping
    public ResponseEntity<Boolean> delete(@RequestParam @Min(value = 1,message = "id сессии должен быть больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.delete(id));
    }

    @PostMapping("/toFinish")
    public ResponseEntity<Boolean> toFinishSession(@RequestParam @Min(value = 1,message = "id сессии должен быть больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.toFinish(id));
    }


    @PatchMapping("/update")
    public ResponseEntity<ResponseSession> update(@Valid @RequestBody RequestUpdateSession requestUpdateSession){
        return ResponseEntity.status(HttpStatus.OK).body(sessionService.update(requestUpdateSession));
    }
}
