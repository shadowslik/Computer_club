package ru.kolobanov.pc.club.entity;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import lombok.*;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import ru.kolobanov.pc.club.controllers.dto.ResponseComputerType;
import ru.kolobanov.pc.club.exeptions.UserNotFoundException;

import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "users")
@Data
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;

    private String number;

    private String email;

    private String password;

    private Double balance = 0.0;

    private Double hours = 0.0;

    private Long token;

    @OneToMany(mappedBy = "user", cascade = CascadeType.REMOVE)
    @JsonIgnore
    private List<Session> sessions = new ArrayList<>();
}
