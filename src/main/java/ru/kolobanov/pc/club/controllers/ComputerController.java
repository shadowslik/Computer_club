package ru.kolobanov.pc.club.controllers;

import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.constraints.Min;
import lombok.AllArgsConstructor;
import org.springframework.validation.annotation.Validated;
import ru.kolobanov.pc.club.controllers.dto.ResponseComputer;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.services.ComputerService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@Tag(name = "Управление компьютерами")
@RequestMapping("api/computers")
@AllArgsConstructor
@Validated
public class ComputerController {

    private ComputerService computerService;

    @PostMapping
    public ResponseEntity<ResponseComputer> addComputer(@RequestParam @Min(value = 1,message = "id компьютера должен быть больше 0") Long type_id){
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(computerService.addComputer(type_id));
    }

    @GetMapping
    public ResponseEntity<List<ResponseComputer>> getComputers(@RequestParam int page,
                                                               @RequestParam int size){
        return ResponseEntity.status(HttpStatus.OK).body(computerService.getComputers(page,size)
                .map(computer -> DtoMapper.computerToResponse(computer)).getContent());
    }

    @GetMapping("/type_id")
    public ResponseEntity<List<ResponseComputer>> getComputersByType_id(@RequestParam int page,
                                                                        @RequestParam int size,
                                                                        @RequestParam @Min(value = 1,message = "id компьютера должен быть больше 0")
                                                                            Long type_id){
        return ResponseEntity.status(HttpStatus.OK).body(computerService.getComputersByType(page, size, type_id)
                .map(computer -> DtoMapper.computerToResponse(computer))
                .getContent());
    }


    @GetMapping("/id")
    public ResponseEntity<ResponseComputer> getComputerById(@RequestParam @Min(value = 1,message = "id компьютера должен быть больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(computerService.getComputerById(id));
    }

    @DeleteMapping
    public ResponseEntity<Boolean> deleteComputer(@RequestParam @Min(value = 1,message = "id компьютера должен быть больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(computerService.deleteComputerById(id));
    }

    @PatchMapping
    public ResponseEntity<ResponseComputer> updateComputer(@RequestParam @Min(value = 1,message = "id компьютера должен быть больше 0") Long id,
                                                           @RequestParam @Min(value = 1,message = "id тарифа должен быть больше 0") Long type_id){
        return ResponseEntity.status(HttpStatus.OK).body(computerService.updateComputer(id,type_id));
    }



}
