package ru.kolobanov.pc.club.controllers;

import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import ru.kolobanov.pc.club.controllers.dto.RequestCreatePromotionDto;
import ru.kolobanov.pc.club.controllers.dto.RequestUpdatePromotionDto;
import ru.kolobanov.pc.club.controllers.dto.ResponsePromotionDto;
import ru.kolobanov.pc.club.services.PromotionService;

import java.util.List;

@RestController
@RequestMapping("api/promotions")
@Tag(name = "Управление акциями")
@AllArgsConstructor
@Validated
public class PromotionController {

    private PromotionService promotionService;

    @GetMapping
    public ResponseEntity<List<ResponsePromotionDto>> getPromotions(){
        return ResponseEntity.status(HttpStatus.OK).body(promotionService.getAllPromotions());
    }

    @PostMapping
    public ResponseEntity<ResponsePromotionDto> postPromotion(@Valid @RequestBody RequestCreatePromotionDto requestCreatePromotionDto){
        return ResponseEntity.status(HttpStatus.CREATED).body(promotionService.postPromotion(requestCreatePromotionDto));
    }

    @PatchMapping
    public ResponseEntity<ResponsePromotionDto> getPromotions(@Valid @RequestBody RequestUpdatePromotionDto  requestUpdatePromotionDto){
        return ResponseEntity.status(HttpStatus.OK).body(promotionService.updatePromotion(requestUpdatePromotionDto));
    }

    @DeleteMapping
    public ResponseEntity<Boolean> getPromotions(@RequestParam @Positive(message = "id акции только больше 0") Long id){
        return ResponseEntity.status(HttpStatus.OK).body(promotionService.deletePromotion(id));
    }
}
