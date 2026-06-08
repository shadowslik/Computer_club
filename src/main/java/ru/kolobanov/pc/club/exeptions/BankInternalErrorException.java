package ru.kolobanov.pc.club.exeptions;

public class BankInternalErrorException extends RuntimeException {
    public BankInternalErrorException(String message) {
        super(message);
    }
}
