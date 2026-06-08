package ru.kolobanov.pc.club.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;

import ru.kolobanov.pc.club.controllers.dto.*;
import ru.kolobanov.pc.club.entity.*;

import ru.kolobanov.pc.club.exeptions.*;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.repository.ComputerRepo;
import ru.kolobanov.pc.club.repository.SessionRepo;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.List;
import java.util.stream.Stream;

@Service
@AllArgsConstructor
public class SessionService {

    SessionRepo sessionRepo;
    ComputerRepo computerRepo;
    UserService userService;
    ComputerService computerService;

    public ResponseSessionAdmin getSession(Long id){
        Session session = checkSessionById(id);
        return DtoMapper.sessionToResponseAdmin(session);
    }

    public Page<ResponseSessionAdmin> getSessions(int page, int size){
        return sessionRepo.findAll(PageRequest.of(page,size))
                .map(DtoMapper::sessionToResponseAdmin);
    }

    public Page<ResponseSessionAdmin> getSessionsByComputerId(Long id, int page, int size){
        computerService.checkById(id);
        return sessionRepo.findByComputer(id,PageRequest.of(page,size))
                .map(DtoMapper::sessionToResponseAdmin);
    }

    public Page<ResponseSession> getSessionsByUserId(Long id, int page, int size){
        userService.checkUserId(id);
        LocalDateTime now = LocalDateTime.now();
        return sessionRepo.findByUser(id,now,PageRequest.of(page,size))
                .map(DtoMapper::sessionToResponse);
    }


    public List<ResponseLastUserSessions> getLastUserSessions(Long id){
        userService.checkUserId(id);
        return sessionRepo.findTop4ByUser_IdOrderByStartTimeDesc(id).stream().map(DtoMapper::sessionToResponseLastUserSessions).toList();

    }

    @Transactional
    public ResponseSession setComputerSession(RequestCreateSession requestCreateSession){

        Long computer_id = requestCreateSession.getComputer_id();
        Long user_id = requestCreateSession.getUser_id();
        Integer hours = requestCreateSession.getDurationHours();
        LocalDateTime start = requestCreateSession.getDateTime();

        User user = userService.checkUserId(user_id);

        LocalDateTime end = start.plusHours(hours);

        Computer computer = computerService.checkById(computer_id);

        checkComputerAvailableForCreate(computer,start,end);

        Double total = computerService.getTotal(hours,computer);

        userService.chekUserBalance(user,total);

        userService.updateBalance(user.getId(),user.getBalance() - total);
        userService.updateHours(user.getId(),user.getHours() + hours);


        Session session = new Session();

        session.setTotal(total);
        session.setComputer(computer);
        session.setUser(user);
        session.setStartTime(start);
        session.setEndTime(end);

        computerRepo.save(computer);

        sessionRepo.save(session);

        releaseComputerIfExpired(computer);

        return DtoMapper.sessionToResponse(session);

    }


    @Transactional
    public ResponseSession update(RequestUpdateSession requestUpdateSession){

        Session session = checkSessionById(requestUpdateSession.getId());

        Computer oldComputer = session.getComputer();
        Computer computer = oldComputer;
        LocalDateTime startTime = session.getStartTime();
        double hours = ChronoUnit.HOURS.between(
                startTime, session.getEndTime()
        );

        if(requestUpdateSession.getComputer_id() != null){
            computer = computerService.checkById(requestUpdateSession.getComputer_id());
        }

        if(requestUpdateSession.getStartTime() != null){
            startTime = requestUpdateSession.getStartTime();
        }

        if(requestUpdateSession.getHours() != null){
            hours = requestUpdateSession.getHours();
        }

        LocalDateTime endTime = startTime.plusHours((long)hours);

        checkComputerAvailableForUpdate(computer, startTime, endTime, session.getId());

        User user = returnUserBalance(session);

        Double total = computerService.getTotal((int)hours,computer);

        userService.chekUserBalance(user,total);

        session.setComputer(computer);
        session.setStartTime(startTime);
        session.setEndTime(endTime);
        session.setTotal(total);

        userService.updateBalance(user.getId(),user.getBalance() - total);
        userService.updateHours(user.getId(),(user.getHours() + hours));

        sessionRepo.save(session);

        if(!computer.equals(oldComputer)){
            releaseComputerIfExpired(computer);
            releaseComputerIfExpired(oldComputer);
        }

        return DtoMapper.sessionToResponse(session);

    }

    @Transactional
    public boolean delete(Long id){
        Session session = checkSessionById(id);
        User user = returnUserBalance(session);
        userService.updateHours(user.getId(), user.getHours());
        userService.updateBalance(user.getId(), user.getBalance());
        sessionRepo.deleteById(id);

        return true;
    }

    @Transactional
    public Boolean toFinish(Long id){

        Session session = checkSessionById(id);

        LocalDateTime now = LocalDateTime.now();

        if(now.isBefore(session.getStartTime())){
            throw new SessionNotStartedException("Нельзя завершить еще не начавшуюся сессию");
        }

        User user = returnUserBalance(session);
        userService.updateHours(user.getId(), user.getHours());
        userService.updateBalance(user.getId(), user.getBalance());

        session.setEndTime(LocalDateTime.now());
        sessionRepo.save(session);

        return true;
    }


    private void releaseComputerIfExpired(Computer computer) {
        LocalDateTime now = LocalDateTime.now();
        Session activeSession = sessionRepo.findCurrentSession(computer, now);
        if (activeSession != null && computer.getStatus() == ComputerStatus.FREE) {
            computer.setStatus(ComputerStatus.BUSY);
            computerRepo.save(computer);
        }else if (activeSession == null && computer.getStatus() == ComputerStatus.BUSY){
            computer.setStatus(ComputerStatus.FREE);
            computerRepo.save(computer);
        }
    }

    private void checkComputerAvailableForCreate(Computer computer, LocalDateTime start, LocalDateTime end) {
        if(sessionRepo.isComputerBusyForCreate(computer,start,end)){
            throw new ComputerBusyException("Этот компьютер будет занят в данное время");
        }
    }

    private void checkComputerAvailableForUpdate(Computer computer, LocalDateTime start, LocalDateTime end, Long id_session) {
        if(sessionRepo.isComputerBusyForUpdate(computer,start,end,id_session)){
            throw new ComputerBusyException("Этот компьютер будет занят в данное время");
        }
    }

    @Transactional
    @Scheduled(fixedDelay = 60000)
    public void releaseAllComputers(){
        try(Stream<Computer> computerList = computerRepo.streamAllBy()){
            computerList.forEach(this::releaseComputerIfExpired);
        }
    }


    private User returnUserBalance(Session session){
        User user = session.getUser();
        LocalDateTime now_time = LocalDateTime.now();
        LocalDateTime start = session.getStartTime();
        LocalDateTime end = session.getEndTime();
        double total_hours = ChronoUnit.HOURS.between(start, end);
        Double price = session.getComputer().getType().getPricePerHour();

        if(now_time.isBefore(start) || now_time.equals(start)){
            user.setBalance(user.getBalance() + total_hours * price);
            user.setHours(user.getHours() - total_hours);
            return user;
        }

        if (now_time.isBefore(end) && now_time.isAfter(start)){

            double lastMinutes = ChronoUnit.MINUTES.between(now_time, end);
            double lastHours = lastMinutes/60;
            user.setBalance(user.getBalance() + lastMinutes * (price / 60));
            user.setHours(user.getHours() - lastHours);
            return user;
        }

        throw new SessionAlreadyFinishedException("Нельзя удалять или отменять уже завершённую сессию");
    }

    private Session checkSessionById(Long id) {
        return sessionRepo.findById(id)
                .orElseThrow(() -> new SessionNotFoundException("Такой сессиии не существует"));
    }




}
