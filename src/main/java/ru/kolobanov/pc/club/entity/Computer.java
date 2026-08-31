package ru.kolobanov.pc.club.entity;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import lombok.Data;

import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "computers")
@Data
public class Computer {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne
    @JoinColumn(name = "type_id")
    private ComputerType type;

    @Enumerated(EnumType.STRING)
    private ComputerStatus status;

    @OneToMany(mappedBy = "computer", cascade = CascadeType.REMOVE)
    @JsonIgnore
    private List<Session> sessions = new ArrayList<>();
}
