# GoMetrum completion, loaded automatically by bash-completion.
_gometrum()
{
    local cur prev words cword split
    _init_completion -s || return

    case "$prev" in
        --config|-config|-c|--service-path|-service-path)
            _filedir
            return
            ;;
    esac

    if $split; then
        case "$prev" in
            --once|--dry-run|--validate|--generate-config|--generate-service|--purge|--version|-v|\
            -once|-dry-run|-validate|-generate-config|-generate-service|-purge|-version)
                COMPREPLY=( $(compgen -W 'true false' -- "$cur") )
                ;;
        esac
        return
    fi

    COMPREPLY=( $(compgen -W '--config -c --once --dry-run --validate
        --generate-config --service-path --generate-service --purge
        --version -v --help -h' -- "$cur") )
}

complete -F _gometrum gometrum
