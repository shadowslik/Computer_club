package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;
import ru.kolobanov.pc.club.entity.ComputerStatus;

@Data
public class ResponseComputer {

    private Long id;

    private String type;

    private Double pricePerHour;

    private ComputerStatus status;
}
