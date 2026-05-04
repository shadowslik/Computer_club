package org.example.entity;

import jakarta.persistence.*;
import lombok.Data;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;

@Entity
@Table(name = "computer_sessions")
@Data
public class ComputerSession{

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne
    @JoinColumn(name = "user_id")
    private User user;

    @ManyToOne
    @JoinColumn(name = "computer_id")
    private Computer computer;

    private OffsetDateTime startTime;

    private OffsetDateTime endTime;

    private Double total;
}
