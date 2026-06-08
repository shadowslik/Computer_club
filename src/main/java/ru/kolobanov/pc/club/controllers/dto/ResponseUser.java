package ru.kolobanov.pc.club.controllers.dto;

import jakarta.persistence.criteria.CriteriaBuilder;
import lombok.Data;

import java.util.List;

@Data
public class ResponseUser {
    private Long id;
    private String name;
    private String email;
    private String phone;
    private Double balance;
    private Double hours;
    private Integer computerSessions;
}
