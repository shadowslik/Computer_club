package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;

@Data
public class ResponseUserTopByBalance {

    private Long id;
    private String name;
    private Double balance;
    private Integer sessions;
}
