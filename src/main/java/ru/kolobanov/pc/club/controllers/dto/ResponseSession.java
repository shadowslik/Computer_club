package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;

import java.time.LocalDateTime;

@Data
public class ResponseSession {

    private Long id;
    private Long computer_id;
    private String computerType;
    private LocalDateTime start;
    private LocalDateTime end;
    private Double total;
}
