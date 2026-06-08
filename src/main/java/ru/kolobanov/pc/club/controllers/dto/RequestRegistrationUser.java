package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

@Data
public class RequestRegistrationUser {

    @NotBlank(message = "Имя обязательно")
    @Size(min = 2, max = 30, message = "Имя должно быть от 2 до 30 символов")
    private String name;

    @NotBlank(message = "Email обязателен")
    @Email(message = "Некоректный формат email")
    private String email;

    @NotBlank(message = "Телефон обязателен")
    @Pattern(regexp = "^\\+?\\d{10,15}$", message = "Неверный формат телефона")
    private String phone;

    @NotBlank(message = "Пароль обязателен")
    @Size(min = 6, message = "Пароль должен быть от 6 символов")
    private String password;
}
