package ru.kolobanov.pc.club.controllers.dto;

import lombok.Data;

import java.time.LocalDateTime;

@Data
public class ResponseSessionAdmin {

    private Long id;
    private Long userId;
    private String userName;
    private String userPhone;
    private String userEmail;
    private Long computer_id;
    private String computerType;
    private LocalDateTime start;
    private LocalDateTime end;
    private Double total;

}
