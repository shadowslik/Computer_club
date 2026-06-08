package ru.kolobanov.pc.club.repository;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import ru.kolobanov.pc.club.entity.Computer;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.stream.Stream;

@Repository
public interface ComputerRepo extends JpaRepository<Computer,Long> {
    Stream<Computer> streamAllBy();

    Page<Computer> findAllByType_idOrderById(Long type_id, Pageable pageable);

    Page<Computer> findAllByOrderById(Pageable pageable);
}
