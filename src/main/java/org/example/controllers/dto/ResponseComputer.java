package org.example.controllers.dto;

import lombok.Data;

@Data
public class ResponseComputer {

    private Long id;

    private String type;

    private String status;

    private Double pricePerHour;
}
