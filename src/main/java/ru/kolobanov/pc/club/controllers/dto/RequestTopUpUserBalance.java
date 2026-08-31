package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.*;
import lombok.Data;

@Data
public class RequestTopUpUserBalance {

    @NotBlank(message = "Номер карты отправителя обязателен")
    @Pattern(regexp = "\\d{16}", message = "Номер карты должен содержать 16 цифр")
    private String cardNumber;

    @NotBlank(message = "CVV карты отправителя обязателен")
    @Pattern(regexp = "\\d{3}", message = "CVV должен содержать 3 цифры")
    private String cardCvv;

    @NotBlank(message = "Срок действия карты отправителя обязателен")
    @Pattern(regexp = "(0[1-9]|1[0-2])/\\d{2}", message = "Формат срока действия: ММ/ГГ")
    private String cardPeriod;

    @NotNull(message = "Сумма не может быть null")
    @Positive(message = "Сумма должна быть положительной")
    private Double amount;
}
