package org.example.controllers.dto;

import lombok.Data;
import org.example.entity.Computer;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;

@Data
public class ResponseComputerSession {

    private Long id;
    private String userName;
    private String userPhone;
    private Computer computer;
    private OffsetDateTime start;
    private OffsetDateTime end;
    private Double total;
}
