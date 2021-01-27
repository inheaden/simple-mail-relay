FROM golang:1.12.0-alpine3.9

ENV USER=appuser
ENV UID=2001
ENV GID=2001

RUN addgroup \
  --gid "${GID}" \
  "${USER}"

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/home/${USER}" \
  --ingroup "${USER}" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

ADD . /home/${USER}

WORKDIR /home/${USER}

RUN go build -o main .
CMD [""]
USER ${USER}

ENTRYPOINT [""]


