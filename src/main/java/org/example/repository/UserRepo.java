package org.example.repository;

import org.example.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

@Repository
public interface UserRepo extends JpaRepository<User,Long> {

    Page<User> findAll(Pageable pageable);


    @Query("SELECT u FROM User u LEFT JOIN u.computerSessions cs GROUP BY u ORDER BY COUNT(cs) DESC")
    Page<User> findAllOrderByComputerSessionsCount(Pageable pageable);

    @Query("SELECT u FROM User u ORDER BY u.hours")
    Page<User> findAllOrderByTotalHours(Pageable pageable);

    @Query("SELECT u FROM User u ORDER BY u.balance DESC")
    Page<User> findAllOrderByBalance(Pageable pageable);

    User findUserByEmail(String email);
    User findUserByNumber(String number);
}
