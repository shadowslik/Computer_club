package org.example.entity;

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

    private String type;

    private String status;

    private Double pricePerHour;

    @OneToMany(mappedBy = "computer", cascade = CascadeType.REMOVE)
    @JsonIgnore
    private List<ComputerSession> computerSessionsId = new ArrayList<>();

    public void addSession(ComputerSession computerSession){
        this.computerSessionsId.add(computerSession);
    }
}
