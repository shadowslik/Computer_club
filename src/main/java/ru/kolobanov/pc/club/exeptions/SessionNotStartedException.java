package ru.kolobanov.pc.club.exeptions;

public class SessionNotStartedException extends RuntimeException {
    public SessionNotStartedException(String message) {
        super(message);
    }
}
