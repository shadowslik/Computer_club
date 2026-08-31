package ru.kolobanov.pc.club.controllers.dto;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;


@Data
public class RequestCreateReferralDto {

    @NotEmpty
    @Min(value = 1, message = "id отравителя не может быть меньше 1")
    private Long idSender;

    @NotEmpty
    @Min(value = 1, message = "id реферала не может быть меньше 1")
    private Long idRef;

}
