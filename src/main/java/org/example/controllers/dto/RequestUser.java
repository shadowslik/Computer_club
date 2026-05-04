package org.example.controllers.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Builder;
import lombok.Data;
import org.springframework.web.bind.annotation.RequestParam;

@Data
public class RequestUser {

    @Schema(nullable = true, example = "null")
    private String name = null;

    @Schema(nullable = true, example = "null")
    private String email = null;

    @Schema(nullable = true, example = "null")
    private String phone = null;

    @Schema(nullable = true, example = "null")
    private String password = null;
}
