package org.example.repository;

import org.example.entity.Computer;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.stream.Stream;

@Repository
public interface ComputerRepo extends JpaRepository<Computer,Long> {
    Stream<Computer> streamAllBy();
}
