package org.example.controllers.dto;

import lombok.Data;
import org.example.entity.ComputerSession;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;

@Data
public class ResponseUser {
    private Long id;
    private String name;
    private String email;
    private String phone;
    private Double balance;
    private Integer hours;
    private List<ResponseComputerSession> computerSessions;
}
