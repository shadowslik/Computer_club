package ru.kolobanov.pc.club.entity;


import jakarta.persistence.*;
import lombok.Data;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

@Entity
@Table(name = "computer_type")
@Data
public class ComputerType {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;

    @Column(name = "price_per_hour")
    private Double pricePerHour;

    private List<String> description = new ArrayList<>();

    public void addDescriptionPar(String... par){
        description.addAll(Arrays.asList(par));
    }
}
