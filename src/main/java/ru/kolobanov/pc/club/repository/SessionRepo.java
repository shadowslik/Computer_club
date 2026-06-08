package ru.kolobanov.pc.club.repository;

import ru.kolobanov.pc.club.entity.Computer;
import ru.kolobanov.pc.club.entity.Session;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;

@Repository
public interface SessionRepo extends JpaRepository<Session,Long> {

    Page<Session> findAll(Pageable pageable);

    @Query("SELECT s FROM Session s " +
            "WHERE s.computer = :computer " +
            "AND s.startTime <= :now " +
            "AND s.endTime > :now")
    Session findCurrentSession(
            @Param("computer") Computer computer,
            @Param("now") LocalDateTime now
    );

    @Query("SELECT CASE WHEN COUNT(s) > 0 THEN true ELSE false END " +
            "FROM Session s " +
            "WHERE s.computer = :computer " +
            "AND s.startTime < :end " +
            "AND s.endTime > :start")
    boolean isComputerBusyForCreate(
            @Param("computer") Computer computer,
            @Param("start") LocalDateTime start,
            @Param("end") LocalDateTime end
    );

    @Query("SELECT CASE WHEN COUNT(s) > 0 THEN true ELSE false END " +
            "FROM Session s " +
            "WHERE s.computer = :computer " +
            "AND s.startTime < :end " +
            "AND s.endTime > :start " +
            "AND s.id <> :id_session")
    boolean isComputerBusyForUpdate(
            @Param("computer") Computer computer,
            @Param("start") LocalDateTime start,
            @Param("end") LocalDateTime end,
            @Param("id_session") Long id_session
    );


    @Query("SELECT s FROM Session s " +
            "WHERE s.computer.id = :computer_id " +
            "ORDER by s.startTime ASC")
    Page<Session> findByComputer(@Param("computer_id") Long id,
                                 Pageable pageable);


    @Query("SELECT s FROM Session s " +
            "WHERE s.user.id = :user_id" +
            " AND s.endTime >= :now_time " +
            "ORDER by s.startTime ASC")
    Page<Session> findByUser(@Param("user_id") Long id,
                             @Param("now_time") LocalDateTime time,
                             Pageable pageable);

    List<Session> findTop4ByUser_IdOrderByStartTimeDesc(Long userId);
}
