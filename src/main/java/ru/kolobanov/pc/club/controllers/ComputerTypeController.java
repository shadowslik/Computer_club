package ru.kolobanov.pc.club.controllers;


import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Min;
import lombok.AllArgsConstructor;
import org.springframework.context.annotation.Scope;
import org.springframework.validation.annotation.Validated;
import ru.kolobanov.pc.club.controllers.dto.RequestCreateComputerType;
import ru.kolobanov.pc.club.controllers.dto.RequestUpdateComputerType;
import ru.kolobanov.pc.club.controllers.dto.ResponseComputerType;
import ru.kolobanov.pc.club.services.ComputerTypeService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@AllArgsConstructor
@RequestMapping("api/computer_types")
@Tag(name = "Управление тарифами")
@Validated
public class ComputerTypeController {

    private ComputerTypeService computerTypeService;


    @GetMapping
    public ResponseEntity<List<ResponseComputerType>> getAllComputerTypes(@RequestParam int page, @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(computerTypeService.getComputerTypes(page,size).getContent());
    }

    @PostMapping
    public ResponseEntity<ResponseComputerType> postComputerType(@Valid @RequestBody RequestCreateComputerType requestCreateComputerType){
        return ResponseEntity.status(HttpStatus.CREATED).body(computerTypeService.addComputerType(requestCreateComputerType));
    }

    @DeleteMapping
    public ResponseEntity<Boolean> deleteComputerType(@RequestParam @Min(value = 1,message = "id тарифа должен быть больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(computerTypeService.deleteTypeById(id));
    }

    @PatchMapping
    public ResponseEntity<ResponseComputerType> updateComputerType(@Valid @RequestBody RequestUpdateComputerType requestUpdateComputerType){
        return ResponseEntity.status(HttpStatus.OK).body(computerTypeService.update(requestUpdateComputerType));
    }
}
