package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;

import java.util.List;

@Data
public class ResponsePromotionDto {

    private Long id;
    private Double value;
    private String type_name;
    private List<String> description;
    private Double beforeValue;
}
