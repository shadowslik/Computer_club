package ru.kolobanov.pc.club.controllers.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.FutureOrPresent;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.Data;

import java.time.LocalDateTime;


@Data
public class RequestCreateSession {

    @NotNull(message = "id компьютера обязателен")
    @Min(value = 1, message = "id компьютера должен быть больше 0")
    private Long computer_id;

    @NotNull(message = "id пользователя обязателен")
    @Min(value = 1, message = "id пользователя должен быть больше 0")
    private Long user_id;

    @NotNull(message = "Количество часов обязательно")
    @Positive(message = "Количество часов должно быть положительным")
    private Integer durationHours;

    @NotNull(message = "Дата обязательна")
    @FutureOrPresent(message = "Дата бронирования не может быть в прошлом")
    @Schema(format = "local-date-time", example = "2026-05-07T19:19:49.324029")
    private LocalDateTime dateTime;
}
