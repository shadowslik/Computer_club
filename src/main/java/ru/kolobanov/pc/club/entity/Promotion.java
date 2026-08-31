package ru.kolobanov.pc.club.entity;


import jakarta.persistence.*;
import lombok.Data;

import java.util.List;

@Entity
@Table(name = "Promotions")
@Data
public class Promotion {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private Double value;

    private Double beforeValue;

    @ElementCollection
    @CollectionTable(
            name = "promotion_descriptions",
            joinColumns = @JoinColumn(name = "promotion_id")
    )
    @Column(name = "description")
    private List<String> description;

    @ManyToOne
    @JoinColumn(name = "type_id")
    private ComputerType type;
}
