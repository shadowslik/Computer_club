package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

@Data
public class RequestUpdateUser {

    @Size(min = 2, max = 30, message = "Имя должно быть от 2 до 30 символов")
    private String name;

    @Email(message = "Некоректный формат email")
    private String email;

    @Pattern(regexp = "^\\+?\\d{10,15}$", message = "Неверный формат телефона")
    private String phone;

    @Size(min = 6, message = "Пароль должен быть от 6 символов")
    private String password;
}
