package org.example.controllers;

import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.AllArgsConstructor;
import org.example.controllers.dto.RequestComputer;
import org.example.controllers.dto.ResponseComputer;
import org.example.entity.Computer;
import org.example.mapper.DtoMapper;
import org.example.services.ComputerService;
import org.springframework.data.domain.Page;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@Tag(name = "Управление компьютерами")
@RequestMapping("api/computers")
@CrossOrigin(origins = "http://localhost:3000")
@AllArgsConstructor
public class ComputerController {

    private ComputerService computerService;

    @PostMapping
    public ResponseEntity<ResponseComputer> addComputer(@RequestParam String type){
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(computerService.addComputer(type));
    }

    @GetMapping
    public ResponseEntity<List<ResponseComputer>> getComputers(){
        List<ResponseComputer> responseComputers = computerService.getComputers().stream()
                .map(computer -> DtoMapper.computerToResponse(computer))
                .toList();
        return ResponseEntity.status(HttpStatus.OK).body(responseComputers);
    }

    @GetMapping("/id/{id}")
    public ResponseEntity<ResponseComputer> getComputerById(@PathVariable Long id){
        return ResponseEntity.ok(computerService.getComputerById(id));
    }

    @DeleteMapping
    public void deleteComputer(@RequestParam Long id){
        computerService.deleteComputerById(id);
    }


    @PutMapping
    public ResponseEntity<ResponseComputer> updateComputer(@RequestBody RequestComputer requestComputer){
        return ResponseEntity.status(HttpStatus.OK)
                .body(computerService.updateComputer(requestComputer));
    }



}
