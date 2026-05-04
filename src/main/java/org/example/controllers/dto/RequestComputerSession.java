package org.example.controllers.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Data;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;


@Data
public class RequestComputerSession {

    @Schema(nullable = true, example = "null", format = "long")
    private Long computer_id;

    @Schema(nullable = true, format = "long")
    private Long user_id;

    @Schema(example = "2", format = "integer")
    private Integer durationHours;
}
