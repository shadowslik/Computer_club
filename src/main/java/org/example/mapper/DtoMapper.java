package org.example.mapper;
import org.example.controllers.dto.*;
import org.example.entity.Computer;
import org.example.entity.ComputerSession;
import org.example.entity.User;

import java.util.Comparator;

public class DtoMapper {

    public static ResponseUser userToResponseUser(User user) {
        ResponseUser responseUser = new ResponseUser();

        responseUser.setId(user.getId());
        responseUser.setName(user.getName());
        responseUser.setPhone(user.getNumber());
        responseUser.setEmail(user.getEmail());
        responseUser.setBalance(user.getBalance());
        responseUser.setComputerSessions(
                        user.getComputerSessions()
                                .stream()
                                .map(DtoMapper::computerSessionToResponse)
                                .toList()
        );
        responseUser.setHours(user.getHours());

        return responseUser;
    }

    public static User requestUserToUser(RequestUser requestUser){
        User user = new User();

        user.setName(requestUser.getName());
        user.setNumber(requestUser.getPhone());
        user.setEmail(requestUser.getEmail());
        user.setPassword(requestUser.getPassword());

        return user;
    }

    public static RequestUser userToRequest(User user ){
        RequestUser requestUser = new RequestUser();

        requestUser.setName(user.getName());
        requestUser.setPhone(user.getNumber());
        requestUser.setEmail(user.getEmail());
        requestUser.setPassword(user.getPassword());

        return requestUser;
    }

    public static ResponseComputerSession computerSessionToResponse(ComputerSession computerSession){

        ResponseComputerSession responseComputerSession = new ResponseComputerSession();

        responseComputerSession.setId(computerSession.getId());
        responseComputerSession.setUserName(computerSession.getUser().getName());
        responseComputerSession.setComputer(computerSession.getComputer());
        responseComputerSession.setUserPhone(computerSession.getUser().getNumber());
        responseComputerSession.setStart(computerSession.getStartTime());
        responseComputerSession.setEnd(computerSession.getEndTime());
        responseComputerSession.setTotal(computerSession.getTotal());

        return responseComputerSession;

    }


    public static ResponseComputer computerToResponse(Computer computer){
        ResponseComputer responseComputer = new ResponseComputer();

        responseComputer.setId(computer.getId());
        responseComputer.setType(computer.getType());
        responseComputer.setPricePerHour(computer.getPricePerHour());
        responseComputer.setStatus(computer.getStatus());

        return responseComputer;
    }

}