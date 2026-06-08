package ru.kolobanov.pc.club.controllers.dto;


import lombok.Data;

import java.util.List;

@Data
public class ResponseComputerType {

    private Long id;
    private String name;
    private Double pricePerHour;
    private List<String> description;
}
