package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

import java.util.List;

@Data
public class RequestCreatePromotionDto {

    @NotNull
    @Min(value = 1, message = "Акция не может быть меньше 1 процента")
    @Max(value = 100, message = "Акция не может быть более 100 процентов")
    private Double value;

    @NotNull
    @Positive(message = "Id типа пк может быть только положительным и не равным 0")
    private Long type_id;

    @NotEmpty(message = "Описание не може бть пустым")
    private List<@NotBlank(message = "Значение описания не может быть пустым") String> description;
}
