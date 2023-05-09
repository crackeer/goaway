function test(params, header, config)
    params['ts'] = os.time()
    local list = sort(params)
    local joinList ={}
    for i, v in ipairs(list) do
        joinList[i] = v[1]..'='..v[2]
    end
    
    local origin = table.concat(joinList, ',')..config['ak']
    params['sign'] = origin
    return header, params
end
    
function sort(data)
    local list = {}
    local index= 1
    for key, val in pairs(data) do
        list[index] = key
        index = index + 1
    end
    table.sort(list, function(a, b) 
        return a < b
    end
    )
    index = 1
    local ret = {}
    for i, value in ipairs(list) do
        local tmp = {}
        tmp[1] = value
        tmp[2] = data[value]
        ret[index] = tmp
        index = index + 1
    end
    return ret
end