package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;

@Data
public class ResponseUserTopBySessions {

    private Long id;
    private String name;
    private Integer sessions;

}
