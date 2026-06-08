package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

@Data
public class RequestLoginUser {

    @NotBlank(message = "Email обязателен")
    @Email(message = "Некоректный формат email")
    private String email;

    @NotBlank(message = "Пароль обязателен")
    private String password;
}
