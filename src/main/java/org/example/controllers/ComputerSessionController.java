package org.example.controllers;


import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.AllArgsConstructor;
import org.example.controllers.dto.RequestComputerSession;
import org.example.controllers.dto.ResponseComputer;
import org.example.controllers.dto.ResponseComputerSession;
import org.example.entity.ComputerSession;
import org.example.services.ComputerSessionService;
import org.springframework.data.domain.Limit;
import org.springframework.data.domain.Page;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@AllArgsConstructor
@RequestMapping("api/computer_sessions")
@CrossOrigin(origins = "http://localhost:3000")
@Tag(name = "Управление сессиями")
public class ComputerSessionController {

    private ComputerSessionService computerSessionService;

    @GetMapping
    public ResponseEntity<List<ComputerSession>> findAll(@RequestParam int page,
                                                         @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(computerSessionService.getAll(page, size).getContent());
    }

    @PostMapping
    public ResponseEntity<ResponseComputerSession> post(@RequestBody RequestComputerSession requestComputerSession){
        return ResponseEntity.status(HttpStatus.CREATED).body(computerSessionService.setComputerSession(requestComputerSession));
    }

    @GetMapping("/id/{id}")
    public ResponseEntity<ResponseComputerSession> getById(@PathVariable Long id){
        return ResponseEntity.status(HttpStatus.OK).body(computerSessionService.getComputerSession(id));
    }

    @GetMapping("/phone/{phone}")
    public ResponseEntity<Page<ResponseComputerSession>> getByPhone(@PathVariable String phone){
        return ResponseEntity.status(HttpStatus.OK).body(computerSessionService.getComputerSessions(phone));
    }

    @GetMapping("/userSessions/{user_id}")
    public ResponseEntity<List<ResponseComputerSession>> getHistory(@PathVariable Long user_id){
        return ResponseEntity.status(HttpStatus.OK).body(computerSessionService.getLastUserSessions(user_id));
    }

    @DeleteMapping
    public ResponseEntity<Boolean> delete(@RequestParam Long id){
        return ResponseEntity.status(HttpStatus.OK).body(computerSessionService.delete(id));
    }


    @PutMapping
    public ResponseEntity<ResponseComputerSession> update(@RequestParam Long id,
            @RequestBody RequestComputerSession requestComputerSession){
        return ResponseEntity.status(HttpStatus.OK).body(computerSessionService.update(id,requestComputerSession));
    }
}
