package ru.kolobanov.pc.club.controllers.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.FutureOrPresent;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.Data;

import java.time.LocalDateTime;

@Data
public class RequestUpdateSession {

    @NotNull(message = "id сессии обязателен")
    @Min(value = 1,message = "id сессии должен быть больше 0")
    private Long id;

    @Min(value = 1,message = "id компьютера должен быть больше 0")
    private Long computer_id;

    @FutureOrPresent(message = "Дата не может быть в прошлом")
    @Schema(format = "local-date-time", example = "2026-05-07T19:19:49.324029")
    private LocalDateTime startTime;

    @Positive(message = "количество часов должно быть больше 0")
    private Integer hours;
}
