package org.example.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;

import org.example.controllers.dto.RequestComputerSession;
import org.example.controllers.dto.ResponseComputer;
import org.example.controllers.dto.ResponseComputerSession;
import org.example.entity.Computer;
import org.example.entity.ComputerSession;
import org.example.entity.User;

import org.example.mapper.DtoMapper;
import org.example.repository.ComputerRepo;
import org.example.repository.ComputerSessionRepo;
import org.example.repository.UserRepo;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
@AllArgsConstructor
public class ComputerSessionService {

    ComputerSessionRepo computerSessionRepo;
    ComputerRepo computerRepo;
    UserService userService;

    public ResponseComputerSession getComputerSession(Long id){
        ComputerSession computerSession = computerSessionRepo.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Такого id нет"));
        return DtoMapper.computerSessionToResponse(computerSession);
    }

    public Page<ResponseComputerSession> getComputerSessions(String phone){
        Pageable pageable = PageRequest.of(0,300);
        Page<ResponseComputerSession> computerSession = computerSessionRepo.findAllByUser_Number(pageable,phone)
                .map(DtoMapper::computerSessionToResponse);

        if(computerSession == null){
            throw new IllegalArgumentException("У пользователя нет такой сессии");
        }
        return computerSession;
    }

    public Page<ComputerSession> getAll(int page, int size){
        return computerSessionRepo.findAll(PageRequest.of(page,size));
    }


    @Transactional
    public List<ResponseComputerSession> getLastUserSessions(Long id){
        User user = userService.checkUserId(id);
        List<ResponseComputerSession> computerSession;

        int lenSessions = user.getComputerSessions().size();

        if(lenSessions > 4){
            computerSession = user.getComputerSessions()
                    .stream()
                    .map(DtoMapper::computerSessionToResponse)
                    .skip(lenSessions - 4)
                    .collect(Collectors.toList());
        }else{
            computerSession = user.getComputerSessions()
                    .stream()
                    .map(DtoMapper::computerSessionToResponse)
                    .collect(Collectors.toList());
        }

        return computerSession;

    }

    @Transactional
    public ResponseComputerSession setComputerSession(RequestComputerSession requestComputerSession){

        Long computer_id = requestComputerSession.getComputer_id();
        Long user_id = requestComputerSession.getUser_id();
        Integer hours = requestComputerSession.getDurationHours();

        if (computer_id == null || user_id == null || hours == null) {
            throw new IllegalArgumentException("Не все поля заполнены (computer_id, user_id, durationHours)");
        }

        if(hours <= 0){
            throw new IllegalArgumentException("Часы указанны не верно");
        }

        Computer computer = checkComputer(computer_id);
        User user = userService.checkUserId(user_id);

        Double total = hours * computer.getPricePerHour();

        releaseComputerIfExpired(computer);

        if(!computer.getStatus().equals("Не занят")){
            throw new IllegalArgumentException("Данный компьютер занят");
        }

        if(total > user.getBalance()){
            throw new IllegalArgumentException("Пользователю надо пополнить баланс");
        }

        user.setBalance(user.getBalance() - total);
        user.setHours(user.getHours() + hours);
        userService.update(user.getId(), DtoMapper.userToRequest(user));

        OffsetDateTime start = OffsetDateTime.now(ZoneOffset.UTC);
        OffsetDateTime end = start.plusHours(hours);

        ComputerSession computerSession = new ComputerSession();

        computerSession.setTotal(total);
        computerSession.setComputer(computer);
        computerSession.setUser(user);
        computerSession.setStartTime(start);
        computerSession.setEndTime(end);

        computer.setStatus("Занят");
        computerRepo.save(computer);

        computerSessionRepo.save(computerSession);

        return DtoMapper.computerSessionToResponse(computerSession);

    }


    @Transactional
    public ResponseComputerSession update(Long id,
            RequestComputerSession requestComputerSession){

        ComputerSession computerSession = computerSessionRepo.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Такой сессиии не существует"));

        if(requestComputerSession.getComputer_id() != null){
            computerSession.setComputer(checkComputer(requestComputerSession.getComputer_id()));
        }

        if(requestComputerSession.getUser_id() != null){
            computerSession.setUser(userService.checkUserId(requestComputerSession.getUser_id()));
        }

        return DtoMapper.computerSessionToResponse(computerSession);

    }

    @Transactional
    public boolean delete(Long id){
        if (!computerSessionRepo.existsById(id)) {
            throw new IllegalArgumentException("Сессия не найдена");
        }
        computerSessionRepo.deleteById(id);
        return true;
    }

    private void releaseComputerIfExpired(Computer computer) {
        ComputerSession activeSession = computerSessionRepo.findByComputerAndEndTimeAfter(
                computer, OffsetDateTime.now(ZoneOffset.UTC));
        if (activeSession == null && computer.getStatus().equals("Занят")) {
            computer.setStatus("Не занят");
            computerRepo.save(computer);
        }
    }

    @Transactional
    @Scheduled(fixedDelay = 60000)
    public void releaseAllComputers(){
        System.out.println("Planner run at " + new Date());
        try(Stream<Computer> computerList = computerRepo.streamAllBy()){
            computerList.forEach(this::releaseComputerIfExpired);
        }
    }


    private Computer checkComputer(Long computer_id){

        Computer computer = computerRepo.findById(computer_id)
                .orElseThrow(() -> new IllegalArgumentException("Компьютера с таким id нету"));

        return computer;
    }

    private boolean chekUserBalance(User user,Double total){
        if(user.getBalance() < total){
            throw new IllegalArgumentException("Пополните пользователю баланс");
        }

        return true;
    }

    private void checkData(OffsetDateTime start, OffsetDateTime end){

        if (start.isBefore(OffsetDateTime.now(ZoneOffset.UTC))){
            throw new IllegalArgumentException("Неправильно указанно время");
        }

        if(start.isAfter(end)){
            throw new IllegalArgumentException("Время начала не может быть меньше конца");

        }

    }

    private Double getTotal(int hours, Computer computer){
        return computer.getPricePerHour() * hours;

    }



}
