# hand-tracking
1. capture camera image
2. hand-tracking 
3. send figure information to some address 

# setup
- `conda env create -f conda-env.yaml `
- `conda activate hand-tracking`

# remove
- `conda deactivate`
- `conda remove -n hand-tracking --all`

# run
- `python main.py {addr} {port} {camera_index}`
    - e.g. `python main.py localhost 8080 1`

<!--
# export 
`conda env export -n hand-tracking`
-->
