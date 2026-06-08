package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

import java.util.List;

@Data
public class RequestCreateComputerType {

    @NotBlank(message = "Имя тарифа обязательно")
    @Size(min = 3,max = 20, message = "Имя тарифа должно быть от 3 до 20 символов")
    private String name;

    @NotNull(message = "Цена тарифа обязательна")
    @Positive(message = "Цена тарифа должна быть положительной")
    private Double pricePerHour;

    @NotEmpty(message = "Описание тарифа обязательно")
    @Size(min = 1, max = 10, message = "В описании должно содержать от 1 до 10 критериев")
    private List<@NotBlank(message = "Критерий описания не может быть пустым") String> description;
}
