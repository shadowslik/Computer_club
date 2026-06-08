package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

import java.util.List;

@Data
public class RequestUpdateComputerType {

    @NotNull(message = "id обязателен")
    @Min(value = 1, message = "id должен быть больше 0")
    private Long id;

    @Size(min = 3,max = 20, message = "Имя тарифа должно быть от 3 до 20 символов")
    private String name;

    @Positive(message = "Цена тарифа должна быть положительной")
    private Double pricePerHour;

    @Size(min = 1, max = 10, message = "Описание должно содержать от 1 до 10 критериев")
    private List<@NotBlank(message = "Критерий описания не может быть пустым") String> description;

}
