#!/bin/bash

CURRDIR=`dirname "$0"`
BASEDIR=`cd "$CURRDIR"; pwd`
export PATH=$PATH:/usr/sbin # the latter dir is not on $PATH for business@203

NAME="wafer-session-server"

CMD="wafer-session-server"

if [ "$1" = "-v" ]; then
    echo "vico_server current version:"
    $BASEDIR/bin/$CMD -v
    exit
fi

if [ "$1" = "-d" ]; then
    shift
    EXECUTEDIR=$1
    shift
else
    EXECUTEDIR=$BASEDIR
fi

if [ ! -d "$EXECUTEDIR" ]; then
    echo "ERROR: $EXECUTEDIR is not a dir"
    exit
fi

if [ ! -d "$EXECUTEDIR"/conf ]; then
    echo "ERROR: could not find $EXECUTEDIR/conf/"
    exit
fi

if [ ! -d "$EXECUTEDIR"/logs ]; then
    mkdir "$EXECUTEDIR"/logs
fi

cd "$EXECUTEDIR"

PID_FILE="$EXECUTEDIR"/logs/"$NAME".pid

check_pid() {
    RETVAL=1
    if [ -f $PID_FILE ]; then
        PID=`cat $PID_FILE`
        cwd_from_pid=`lsof -p $PID 2>/dev/null | awk '/cwd/{ print $9 }'`
        echo "Pid:$PID"
        echo "WorkingDir:$cwd_from_pid"
        if [ "$cwd_from_pid" = "$BASEDIR" ]; then
            echo "service running at dir: $cwd_from_pid with pid:$PID"
            RETVAL=0
        fi
    fi
}

check_running() {
    PID=0
    RETVAL=0
    check_pid
    if [ $RETVAL -eq 0 ]; then
        echo "$CMD is running as $PID, we'll do nothing"
        exit
    fi
}


start() {
    check_running
    echo "starting $CMD ..."
    "$BASEDIR"/"$CMD" -d "$EXECUTEDIR" 2>"$EXECUTEDIR"/logs/"$NAME".err >"$EXECUTEDIR"/logs/"$NAME".out &
    PID=$!
    echo $PID > "$PID_FILE"
        sleep 1
        status
}

stop() {
    check_pid
    if [ $RETVAL -eq 0 ]; then
        echo "$CMD is running as $PID, stopping it..."
        #kill -9 $PID
        kill -15 $PID
                sleep 1
        echo "done"
    else
        echo "$CMD is not running, do nothing"
    fi

        while true; do
                check_pid
                if [ $RETVAL -eq 0 ]; then
                        echo "$CMD is running, waiting it's exit..."
                        sleep 1
                else
                        echo "$CMD is stopped safely, you can restart it now"
                        break
                fi
        done

    if [ -f $PID_FILE ]; then
        rm $PID_FILE
    fi
}

status() {
    check_pid
    STDOUT="$EXECUTEDIR"/logs/"$NAME".out
    if [ $RETVAL -eq 0 ]; then
        echo "$CMD is running as $PID ..."
        OFFSET1=`wc -l <"$STDOUT"`
        kill -USR1 "$PID"
        sleep 0.5
        OFFSET2=`wc -l <"$STDOUT"`
        tail -n $(($OFFSET2 - $OFFSET1)) "$STDOUT"
    else
        echo "$CMD is not running"
    fi
}

RETVAL=0
case "$1" in
    start)
        start $@
        ;;  
    stop)
        stop
        ;;  
    restart)
        stop
        start $@
        ;;  
    status)
        status
        ;;  
    *)  
        echo "Version: $VERSION"
        echo "Usage: $0 [-d EXECUTION_PATH] {start|stop|restart|status}"
        RETVAL=1
esac
exit $RETVAL

