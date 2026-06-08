package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;

@Data
public class ResponseUserTopByHours {

    private Long id;
    private String name;
    private Double hours;
    private Integer sessions;
}
