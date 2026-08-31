package ru.kolobanov.pc.club.mapper;
import ru.kolobanov.pc.club.controllers.dto.*;
import ru.kolobanov.pc.club.entity.*;

import java.time.temporal.ChronoUnit;

public class DtoMapper {

    public static ResponseUser userToResponseUser(User user) {
        ResponseUser responseUser = new ResponseUser();

        responseUser.setId(user.getId());
        responseUser.setName(user.getName());
        responseUser.setPhone(user.getNumber());
        responseUser.setEmail(user.getEmail());
        responseUser.setBalance(user.getBalance());
        responseUser.setComputerSessions(user.getSessions().size());
        responseUser.setHours(user.getHours());
        responseUser.setToken(user.getToken());

        return responseUser;
    }

    public static ResponseUserTopBySessions userToResponseTopBySessions(User user){
        ResponseUserTopBySessions responseUserTopBySessions = new ResponseUserTopBySessions();

        responseUserTopBySessions.setId(user.getId());
        responseUserTopBySessions.setSessions(user.getSessions().size());
        responseUserTopBySessions.setName(user.getName());

        return responseUserTopBySessions;
    }

    public static ResponseUserTopByHours userToResponseTopByHours(User user){
        ResponseUserTopByHours responseUserTopByHours = new ResponseUserTopByHours();

        responseUserTopByHours.setId(user.getId());
        responseUserTopByHours.setSessions(user.getSessions().size());
        responseUserTopByHours.setName(user.getName());
        responseUserTopByHours.setHours(user.getHours());

        return responseUserTopByHours;
    }

    public static ResponseUserTopByBalance userToResponseTopByBalance(User user){
        ResponseUserTopByBalance responseUserTopByBalance = new ResponseUserTopByBalance();

        responseUserTopByBalance.setId(user.getId());
        responseUserTopByBalance.setSessions(user.getSessions().size());
        responseUserTopByBalance.setName(user.getName());
        responseUserTopByBalance.setBalance(user.getBalance());

        return responseUserTopByBalance;
    }

    public static User requestRegistrationUserToUser(RequestRegistrationUser requestRegistrationUser){

        User user = new User();

        user.setName(requestRegistrationUser.getName());
        user.setNumber(requestRegistrationUser.getPhone());
        user.setEmail(requestRegistrationUser.getEmail());
        user.setPassword(requestRegistrationUser.getPassword());

        return user;
    }

    public static ResponseSession sessionToResponse(Session session){

        ResponseSession responseSession = new ResponseSession();

        responseSession.setId(session.getId());
        responseSession.setComputer_id(session.getComputer().getId());
        responseSession.setComputerType(session.getComputer().getType().getName());
        responseSession.setStart(session.getStartTime());
        responseSession.setEnd(session.getEndTime());
        responseSession.setTotal(session.getTotal());

        return responseSession;

    }

    public static ResponseLastUserSessions sessionToResponseLastUserSessions(Session session){

        ResponseLastUserSessions responseLastUserSessions = new ResponseLastUserSessions();

        responseLastUserSessions.setId(session.getId());
        responseLastUserSessions.setComputer_id(session.getComputer().getId());
        responseLastUserSessions.setComputerType(session.getComputer().getType().getName());
        responseLastUserSessions.setStart(session.getStartTime());

        double hours = ChronoUnit.HOURS.between(
                session.getStartTime(), session.getEndTime()
        );

        responseLastUserSessions.setHours(hours);
        responseLastUserSessions.setTotal(session.getTotal());

        return responseLastUserSessions;

    }

    public static ResponseSessionAdmin sessionToResponseAdmin(Session session){

        ResponseSessionAdmin responseSessionAdmin = new ResponseSessionAdmin();

        responseSessionAdmin.setId(session.getId());
        responseSessionAdmin.setUserId(session.getUser().getId());
        responseSessionAdmin.setUserName(session.getUser().getName());
        responseSessionAdmin.setUserEmail(session.getUser().getEmail());
        responseSessionAdmin.setUserPhone(session.getUser().getNumber());
        responseSessionAdmin.setComputer_id(session.getComputer().getId());
        responseSessionAdmin.setComputerType(session.getComputer().getType().getName());
        responseSessionAdmin.setStart(session.getStartTime());
        responseSessionAdmin.setEnd(session.getEndTime());
        responseSessionAdmin.setTotal(session.getTotal());

        return responseSessionAdmin;

    }


    public static ResponseComputer computerToResponse(Computer computer){
        ResponseComputer responseComputer = new ResponseComputer();

        responseComputer.setId(computer.getId());
        responseComputer.setType(computer.getType().getName());
        responseComputer.setPricePerHour(computer.getType().getPricePerHour());
        responseComputer.setStatus(computer.getStatus());

        return responseComputer;
    }


    public static ComputerType requestToComputerType(RequestCreateComputerType requestCreateComputerType){
        ComputerType computerType = new ComputerType();

        computerType.setName(requestCreateComputerType.getName());
        computerType.setPricePerHour(requestCreateComputerType.getPricePerHour());
        computerType.setDescription(requestCreateComputerType.getDescription());

        return computerType;
    }

    public static RequestUpdateComputerType computerTypeToUpdateRequest(ComputerType computerType){

        RequestUpdateComputerType requestUpdateComputerType = new RequestUpdateComputerType();

        requestUpdateComputerType.setId(computerType.getId());
        requestUpdateComputerType.setName(computerType.getName());
        requestUpdateComputerType.setPricePerHour(computerType.getPricePerHour());
        requestUpdateComputerType.setDescription(computerType.getDescription());

        return requestUpdateComputerType;
    }

    public static ResponseComputerType computerTypeToResponse(ComputerType computerType){
        ResponseComputerType responseComputerType = new ResponseComputerType();

        responseComputerType.setId(computerType.getId());
        responseComputerType.setName(computerType.getName());
        responseComputerType.setPricePerHour(computerType.getPricePerHour());
        responseComputerType.setDescription(computerType.getDescription());

        return responseComputerType;
    }

    public static ResponsePromotionDto promotionToResponse(Promotion promotion){
        ResponsePromotionDto responsePromotionDto = new ResponsePromotionDto();

        responsePromotionDto.setDescription(promotion.getDescription());
        responsePromotionDto.setValue(promotion.getValue());
        responsePromotionDto.setType_name(promotion.getType().getName());
        responsePromotionDto.setId(promotion.getId());
        responsePromotionDto.setBeforeValue(promotion.getBeforeValue());

        return responsePromotionDto;
    }


    public static Promotion promotionRequestToPromotion(RequestCreatePromotionDto requestCreatePromotionDto, ComputerType type){

        Promotion promotion = new Promotion();

        promotion.setDescription(requestCreatePromotionDto.getDescription());
        promotion.setValue(requestCreatePromotionDto.getValue());
        promotion.setType(type);

        return promotion;
    }


    public static Referrals referralsRequestToReferral(Long sender, Long ref){

        Referrals referrals = new Referrals();

        referrals.setIdSender(sender);
        referrals.setIdRef(ref);

        return referrals;
    }


}