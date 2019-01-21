package com.cbe.domain;

import javax.persistence.Entity;
import javax.persistence.Id;

@Entity
public class PersonType {

    @Id
    private long id;
    private String type;

    public PersonType() {
    }

    public PersonType(final long id, final String type) {
        this.id = id;
        this.type = type;
    }

    public long getId() {
        return id;
    }

    public void setId(final long id) {
        this.id = id;
    }

    public String getType() {
        return type;
    }

    public void setType(final String type) {
        this.type = type;
    }

    @Override
    public String toString() {
        return "PersonType{" +
                "id=" + id +
                ", type='" + type + '\'' +
                '}';
    }
}
